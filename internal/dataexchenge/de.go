package dataexchenge

import "vpngigabot/internal/models"

type DE struct {
	de chan models.DataExchenge
}

func New() *DE {
	return &DE{de: make(chan models.DataExchenge)}
}

func (de *DE) Read() chan models.DataExchenge {
	return de.de
}

func (de *DE) Write(datEx *models.DataExchenge) {
	de.de <- *datEx
}
