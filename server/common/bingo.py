import logging

from common.utils import Bet, store_bets


def _transformIntoBingoDTO(info: str):
    bingoDTO = {
        'name': '',
        'document': '',
        'born_date': '',
        'number': '',
        'surname': '',
    }
    for kv in info.lstrip('{').rstrip('}').split(','):
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
    print(f"res is: {res}")
    return res


class Bingo:

    def __init__(self, agency: str):
        self.agency = agency

    def processMessage(self, data: str):
        try:
            infoOfBets = _stringToArrayStringBingoDto(data)  # ToDo check here
            bets = []
            for b in infoOfBets:
                bets.append(self._transformIntoBet(_transformIntoBingoDTO(b)))
            store_bets(bets)
            amount = len(bets)
            logging.info(f'action: apuestas_almacenadas | result: success | cantidad: {amount}')
            return amount
        except Exception as e:
            logging.error(f"error happened while processing bet, {e}")
            return 0

    def _transformIntoBet(self, bingoDto: dict):
        return Bet(self.agency, bingoDto['name'], bingoDto['surname'], bingoDto['document'], bingoDto['born_date'],
                   bingoDto['number'])
