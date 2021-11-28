package output

import (
	"encoding/json"
	"fmt"

	"github.com/zartbot/goflow/datarecord"
)

func ConsoleLog(prefix string, dfchan chan *datarecord.DataFrame) {
	for {
		d := <-dfchan
		RecordMap(d)
		d.Print(prefix)

		if (d.Type != "NULL") && (d.Type != "OptionTemplate") {

			fmt.Println("------------")
			r, _ := json.Marshal(d.RecordList())
			fmt.Println(string(r))
			fmt.Println("------------")

		}

	}
}
