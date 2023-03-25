import json
import logging

from common.utils import Bet, store_bets, load_bets, has_won


def _transformIntoBingoDTO(data: dict):
    return {
        "name": data['name'],
        "document": data['document'],
        "born_date": data['born_date'],
        "number": data['number'],
        "surname": data['surname'],
    }


class LotteryManager:
    def __init__(self, lmReading, agenciesFinished, agenciesNeeded: int = 5):
        self._winners = []
        self._agenciesFinished = agenciesFinished
        self._minimumAmountOfAgencies = agenciesNeeded
        self._amountReading = lmReading

    def agencyFinished(self, agency):
        with self._agenciesFinished.get_lock():
            self._agenciesFinished.value += 1

    def getWinners(self, lock):
        if self._agenciesFinished.value < self._minimumAmountOfAgencies:
            raise IndexError('can`t start the lottery, still missing agencies to report its winners')
        with self._amountReading.get_lock():  # This could suffer from starvation to anyone trying to read, but is min
            if self._amountReading == 0:
                lock.acquire()
            self._amountReading.value += 1
        for x in load_bets():
            if has_won(x):
                self._winners.append(str(x.document))
        with self._amountReading.get_lock():
            self._amountReading.value -= 1
            if self._amountReading == 0:
                lock.release()
        return self._winners


class Bingo:
    def __init__(self, agency: str, lock, lm: LotteryManager, agenciesFinished=0):
        self.agency = agency
        self._agencyFinished = agenciesFinished
        self.__actionsMap = {'sendingBatch': self._processBets, 'findMeMyOgre': self._findWinners}
        self._lock = lock
        self._lm = lm

    def processMessage(self, data: str) -> dict:
        action = 'unknown'
        try:
            infoOfBets = json.loads(data)
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
        self._lock.acquire()
        store_bets(bets)
        self._lock.release()
        amount = len(bets)
        logging.info(f'action: apuestas_almacenadas | result: success | cantidad: {amount}')
        return {'amountProcessed': amount}

    def _findWinners(self, *args):
        lotteryManager = self._lm
        try:
            info = {'winners': lotteryManager.getWinners(self._lock), 'status': 'foundOgre'}
            return info
        except IndexError:
            return {'status': 'notAllOgre'}

    def _transformIntoBet(self, bingoDto: dict):
        return Bet(self.agency, bingoDto['name'], bingoDto['surname'], bingoDto['document'], bingoDto['born_date'],
                   bingoDto['number'])
