import json
import logging

from common.utils import Bet, store_bets


def _transformIntoBingoDTO(data: dict):
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
            infoOfBets = json.loads(data)
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
