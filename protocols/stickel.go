package protocols

import (
	"fmt"
	"go_braids/braids"
	"os"
)

const NodeCount = 10

func StikelProtocolImplementation() {
	fmt.Print("Реализация протокола Стикеля.\n\n")
	var width int
	fmt.Print(`Введите "ширину" косы: `)
	_, _ = fmt.Fscan(os.Stdin, &width)
	fmt.Scanln()
	pubBraid1 := braids.GenRandomBraid(width, NodeCount)
	pubBraid2 := braids.GenRandomBraid(width, NodeCount)
	var (
		n, m int //Секреты Алисы
		r, s int //Секреты Боба
	)
	fmt.Print(`Введите "Секреты" Алисы в формате x, y: `)
	_, _ = fmt.Scanf("%d, %d\n", &n, &m)
	fmt.Scan()
	fmt.Print(`Введите "Секреты" Боба в формате x, y: `)
	_, _ = fmt.Scanf("%d, %d\n", &r, &s)
	//Шаг (1) из прикрепленного документа
	aliceComputeAN := braids.Pow(pubBraid1, n)
	aliceComputeBM := braids.Pow(pubBraid2, n)
	aliceSend := braids.Mult(aliceComputeAN, aliceComputeBM)
	fmt.Println("Алиса отправила свою часть ключа")
	//Шаг (2) из прикрепленного документа
	bobComputeAR := braids.Pow(pubBraid1, r)
	bobComputeBS := braids.Pow(pubBraid2, s)
	bobSend := braids.Mult(bobComputeAR, bobComputeBS)
	fmt.Println("Боб отправил свою часть ключа")
	//Шаги 3 и 4 соответственно
	aliceDecipher := braids.Mult(aliceComputeAN, bobSend, aliceComputeBM)
	bobDecipher := braids.Mult(bobComputeAR, aliceSend, bobComputeBS)
	//Эпилог и вывод информации о том, что вообще было обработано
	if aliceDecipher.Equals(bobDecipher) {
		fmt.Println("Обмен сообщениями прошёл успешно! Алиса и Боб установили связь и сформировали общий секрет!")
	} else {
		fmt.Println("Ошибка! Общий секрет не установлен!")
	}
	fmt.Println("В главных ролях:")
	fmt.Println("Первая публичная коса: ", pubBraid1)
	fmt.Println("Вторая публичная коса: ", pubBraid2)
	fmt.Printf("Секреты Алисы: %d, %d \n", n, m)
	fmt.Printf("Секреты Боба: %d, %d \n", r, s)
	fmt.Println("Часть общего ключа (Алиса): \n", aliceSend)
	fmt.Println("Часть общего ключа (Боб): \n", bobSend)
	fmt.Println("Расшифровка Алисы: \n", aliceDecipher)
	fmt.Println("Расшифровка Боба: \n", bobDecipher)

	fmt.Scanln()
}
