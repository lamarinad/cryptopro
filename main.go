package main

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"log"
	"math/big"

	"github.com/lamarinad/cryptopro/internal/exchange"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/fbsobreira/gotron-sdk/pkg/client"
	"google.golang.org/protobuf/proto"
)

func main() {
	privKey := generateAccount()
	addr := "TH1sXGcekx4FbmTo7kmuxfWGBDNBjovpTR"
	fmt.Println(exchange.GetTronBalance(addr))

	sendTrx(privKey, addr, 300)

	ownerAddress := "TH1sXGcekx4FbmTo7kmuxfWGBDNBjovpTR"
	recipientAddress := "адрес гриши"
	tokenName := "Marinad"
	totalSupply := 1000000000

	// Выделение токена на чей-то кошелек
	contractAddress := "адрес вашего контракта" // Надо будет получить после создания контракта
	amount := 1000000
	MND, err = sendToken(ownerAddress, recipientAddress, contractAddress, amount)
	if err != nil {
		fmt.Println("Ошибка выделения токенов:", err)
		return
	}
	fmt.Println("Токены выделены. Transaction ID:", MND)
}

func createToken(ownerAddress, name string, supply int) (string, error) {
	contract := Contract{
		OwnerAddress:    ownerAddress,
		Name:            name,
		TotalSupply:     supply,
		TRC20Compatible: true,
		Precision:       6, // Зависит от ваших требований
	}
}

// генерация приватного ключа
func generateAccount() string {
	privKey, err := ecdsa.GenerateKey(btcec.S256(), rand.Reader)
	if err != nil {
		log.Fatalf("Ошибка генерации приватного ключа: %v", err)
	}

	tronAddress := address.PubkeyToAddress(privKey.PublicKey)

	fmt.Printf("%x\n", privKey.D.Bytes())
	fmt.Println(tronAddress)

	return fmt.Sprintf("%x", privKey.D.Bytes())
}

func sendTrx(privateKey string, toAddress string, amount int64) (string, error) {
	// Подключение к Tron-клиенту
	c := client.NewGrpcClient("grpc.trongrid.io:50051")
	c.Start()
	defer c.Stop()
	// Преобразование приватного ключа из строки
	privKeyBytes := big.NewInt(0)
	privKeyBytes.SetString(privateKey, 16)
	privKey, _ := btcec.PrivKeyFromBytes(privKeyBytes.Bytes())

	// Получение адреса отправителя из приватного ключа
	fromAddress := address.PubkeyToAddress(privKey.ToECDSA().PublicKey).String()

	tx, err := c.Transfer(fromAddress, toAddress, amount)
	if err != nil {
		return "", fmt.Errorf("Ошибка создания транзакции: %v", err)
	}

	// Подпись транзакции приватным ключом
	rawData, err := proto.Marshal(tx.GetTransaction())
	if err != nil {
		return "", fmt.Errorf("Ошибка сериализации транзакции: %v", err)
	}

	hash := sha256.Sum256(rawData)

	// Использование ecdsa для подписи
	r, s, err := ecdsa.Sign(rand.Reader, privKey.ToECDSA(), hash[:])
	if err != nil {
		return "", fmt.Errorf("Ошибка подписи транзакции: %v", err)
	}

	// Объединение r и s в одну подпись
	signature := append(r.Bytes(), s.Bytes()...)

	// Присоединение подписи к транзакции
	tx.Transaction.Signature = append(tx.Transaction.Signature, signature)

	// Отправка транзакции
	result, err := c.Broadcast(tx.Transaction)
	if err != nil {
		return "", fmt.Errorf("Ошибка отправки транзакции: %v", err)
	}

	return result.String(), nil
}
