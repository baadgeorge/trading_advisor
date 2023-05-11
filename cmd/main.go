package main

import (
	"final/telegram"
	TGApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {
	if err := godotenv.Load("passw.env"); err != nil {
		logrus.Fatalf("can't load env: %s", err.Error())
		return
	}

	//rService := database.NewRepoService(db)
	//instr := sdk.NewServicePool("t.LCgR7EtlcOlfRM1epjYTV0Z9bF2gqeGZhUH831L8vUBb0-LH6EuIXEW5o6k4XKmpA7vPVw39hU02de1Mhcv-yw")
	botAPI, err := TGApi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		panic(err)
	}
	log := logrus.New()
	log.SetReportCaller(true)
	bot := telegram.NewBot(botAPI)
	err = bot.StartTelegramUpdates()

	if err != nil {
		logrus.Fatal(err)
		return
	}

}

/*func main() {
	//instr := sdk.NewServicePool()
	//getFigiList(instr.InstrumentsService, shares_type)
	//ctx := context.Background()

	tb := trade.NewTradeBot("t.LCgR7EtlcOlfRM1epjYTV0Z9bF2gqeGZhUH831L8vUBb0-LH6EuIXEW5o6k4XKmpA7vPVw39hU02de1Mhcv-yw")
	ch := tb.GetNewWorkersConfigCh()

	wg := new(sync.WaitGroup)
	wg.Add(3)
	go tb.TradeBotRun(wg)
	go sendStrat(ch, wg)
	time.Sleep(30 * time.Second)
	go stopWorkers(tb, wg)
	wg.Wait()
}*/

/*func stopWorkers(tb *trade.TradeBot, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		for k := range tb.workers {
			time.Sleep(5 * time.Second)
			tb.StopWorker(k)
		}
	}
}

func sendStrat(ch chan *trade.WorkerConfig, wg *sync.WaitGroup) {
	defer wg.Done()
	//dema1 := strategy.NewDoubleEMA(5, 20, 1, 0, 0)
	bb2 := strategy.NewBollingerBand(1, 5, 2)
	//dema2 := strategy.NewDoubleEMA(5, 20, 1, 0, 0)

	workerBB := trade.NewWorkerConfig("BBG000BNNYW1", 20, bb2)
	//workerDEMA1 := trade.NewWorkerConfig("BBG000BNNYW1", 10, dema1)
	//workerDEMA2 := trade.NewWorkerConfig("BBG000BNNYW1", 5, dema2)

	time.Sleep(2 * time.Second)
	ch <- workerBB
	time.Sleep(3 * time.Second)
	//ch <- workerDEMA1
	time.Sleep(5 * time.Second)
	//ch <- workerDEMA2
}
*/
