import logging

from common.utils import Bet, store_bets, load_bets, has_won


def _transformIntoBingoDTO(info: str):
    bingoDTO = {
        'name': '',
        'document': '',
        'born_date': '',
        'number': '',
        'surname': '',
    }
    for kv in info.lstrip(',').lstrip('{').rstrip('}').split(','):
        k, v = kv.split(':')
        if k in bingoDTO:
            bingoDTO[k] = v
    return bingoDTO


def _stringToArrayStringBingoDto(msg: str):
    res = []
    inter = ''
    for x in msg.rstrip(']').lstrip('['):
        if x == '"':
            continue
        inter += x
        if x == '}':
            res.append(inter)
            inter = ''
    return res


def _stringToActionObject(msg: str) -> dict:
    d = {'action': '', 'data': []}
    data = [x for x in msg.split('"') if x != '']
    for i in range(len(data)):
        if data[i] == 'action':
            d['action'] = data[i + 2]  # position of the action after split and filter
            i += 2
        if data[i] == 'data':
            finalPositionOfArray = i+2
            for j in range(i, len(data)):
                if ']' in data[j]:
                    finalPositionOfArray = j + 1
                    break
            d['data'] = _stringToArrayStringBingoDto("".join(data[i+1:finalPositionOfArray]).rstrip("}").lstrip(":"))
            i += finalPositionOfArray
    return d


class Bingo:

    def __init__(self, agency: str, agenciesFinished=0):
        self.agency = agency
        self._agencyFinished = agenciesFinished
        self.__actionsMap = {'sendingBatch': self._processBets, 'findMeMyOgre': self._findWinners}

    def processMessage(self, data: str) -> dict:
        action = 'unknown'
        try:
            infoOfBets = _stringToActionObject(data)
            action = infoOfBets["action"]
            methodToUse = self.__actionsMap.get(action, None)
            if not methodToUse:
                raise Exception('unknown action, please try again with a different action')
            return self.__actionsMap[action](infoOfBets.get('data', []))
        except Exception as e:
            logging.error(f"error happened while processing action: {action}, error is: {e}")
            return {'amountProcessed': 0}

    def _processBets(self, infoOfBets):
        bets = []
        for b in infoOfBets:
            bets.append(self._transformIntoBet(_transformIntoBingoDTO(b)))
        store_bets(bets)
        amount = len(bets)
        logging.info(f'action: apuestas_almacenadas | result: success | cantidad: {amount}')
        return {'amountProcessed': amount}

    def _findWinners(self, *args):
        lotteryManager = LotteryManager()
        try:
            info = {'winners': lotteryManager.getWinners(), 'status': 'foundOgre'}
            return info
        except IndexError:
            return {'status': 'notAllOgre'}

    def _transformIntoBet(self, bingoDto: dict):
        return Bet(self.agency, bingoDto['name'], bingoDto['surname'], bingoDto['document'], bingoDto['born_date'],
                   bingoDto['number'])


class LotteryManager:
    def __int__(self, agenciesNeeded: int = 5):
        self._winners = []
        self._agenciesFinished = set()
        self._minimumAmountOfAgencies = agenciesNeeded

    def __new__(cls, agenciesNeeded: int = 5):
        if not hasattr(cls, 'instance'):
            cls.instance = super(LotteryManager, cls).__new__(cls)
            cls._winners = []
            cls._agenciesFinished = set()
            cls._minimumAmountOfAgencies = agenciesNeeded
        return cls.instance

    def agencyFinished(self, agency):
        self._agenciesFinished.add(agency)

    def getWinners(self):
        if len(self._agenciesFinished) < self._minimumAmountOfAgencies:
            raise IndexError('can`t start the lottery, still missing agencies to report its winners')
        if not self._winners:
            # ToDo later here we will need a sync feature
            for x in load_bets():
                if has_won(x):
                    self._winners.append(str(x.document))
            logging.info('action: sorteo | result: success')
        return self._winners
