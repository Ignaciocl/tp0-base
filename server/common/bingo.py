import json
import logging

from common.utils import Bet, store_bets


def _transformIntoBingoDTO(info: str):
    data = json.loads(info)
    return {
        "name": data['name'],
        "document": data['document'],
        "born_date": data['born_date'],
        "number": data['number'],
        "surname": data['surname'],
    }


class Bingo:

    def __init__(self, agency: str):
        self.agency = agency

    def processMessage(self, data: str):
        try:
            bet = self._transformIntoBet(_transformIntoBingoDTO(data))
            store_bets([bet])
            logging.info(f'action: apuesta_almacenada | result: success | dni: ${bet.document} | numero: ${bet.number}')
        except Exception as e:
            logging.error(f"error happened while processing bet, {e}")
            return

    def _transformIntoBet(self, bingoDto: dict):
        return Bet(self.agency, bingoDto['name'], bingoDto['surname'], bingoDto['document'], bingoDto['born_date'],
                   bingoDto['number'])
