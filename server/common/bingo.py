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
        k, v = kv.strip().split(':')
        k, v = k.strip('"'), v.strip('"')
        if k in bingoDTO:
            bingoDTO[k] = v
    return bingoDTO


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
