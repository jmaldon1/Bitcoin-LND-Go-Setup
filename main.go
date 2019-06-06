package main

import(
    "fmt"
    "os/user"
    "log"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials"
    "google.golang.org/grpc/metadata"
    "time"
    "context"
    "github.com/lightningnetwork/lnd/lnrpc"
)

// Macaroons global var
var MACAROONOPTION grpc.CallOption

func main(){
    grpcConn := grpcSetup()
    defer grpcConn.Close()

    lncli := lnrpc.NewLightningClient(grpcConn)

    ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
    walletBalanceReq := lnrpc.WalletBalanceRequest{}
    walletRes, err := lncli.WalletBalance(ctx, &walletBalanceReq, MACAROONOPTION)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println(walletRes.TotalBalance)

    newAddressRequest := lnrpc.NewAddressRequest{Type: 0}
    newAddrRes, err := lncli.NewAddress(ctx, &newAddressRequest, MACAROONOPTION)
    if err != nil {
        fmt.Println(err)
    }

    fmt.Println(newAddrRes.Address)

}

func grpcSetup()*grpc.ClientConn{
    usr, err := user.Current()
    if err != nil {
        log.Fatal( err )
    }
    homeDir := usr.HomeDir
    lndDir := fmt.Sprintf("%s/Library/Application Support/Lnd", homeDir)

    // SSL credentials setup
    var serverName string
    certFileLocation := fmt.Sprintf("%s/tls.cert", lndDir)
    creds, err := credentials.NewClientTLSFromFile(certFileLocation, serverName)
    if err != nil {
        fmt.Println(err)
    }

    // Macaroon setup
    macaroonFileLocation := fmt.Sprintf("%s/data/chain/bitcoin/regtest/admin.macaroon", lndDir)
    macaroonMap := map[string]string{"macaroon": macaroonFileLocation}
    macaroonMetadata := metadata.New(macaroonMap)
    MACAROONOPTION = grpc.Header(&macaroonMetadata)

    conn, err := grpc.Dial("localhost:10009", grpc.WithTransportCredentials(creds))
    if err != nil {
        log.Fatal( err )
    }

    return conn
}