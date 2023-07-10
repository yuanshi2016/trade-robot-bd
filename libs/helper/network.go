/**
 * @Notes:
 * @class network
 * @package
 * @author: 原始
 * @Time: 2023/6/11   21:03
 */
package helper

import (
	"io"
	"log"
	"net/http"
	"os"
)

func Download(url string, filename string) {
	res, err := http.Get(url)
	if err != nil {
		log.Printf("http.Get -> %v", err.Error())
		return
	}
	data, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("ioutil.ReadAll -> %s", err.Error())
		return
	}
	defer res.Body.Close()

	if err = os.WriteFile(filename, data, 0666); err != nil {
		log.Println("Error Saving:", filename, err.Error())
	} else {
		log.Println("Saved:", filename)
	}

}
