# fptai-sdk-go
FPT.AI SDK for the Go programming language

## fptai CLI
Installation and usage
```
$ export GOPATH=$(pwd)
$ export PATH=$PATH:$GOPATH/bin
$ go get github.com/fpt-corp/fptai-sdk-go/cmd/fptai
$ fptai help
$ fptai train -t intent -i train.csv -u your_username -p your_password -c your_application_code
$ fptai test -t intent -i test.csv -u your_username -p your_password -c your_application_code
```

training.csv and test.csv file must be a CSV file and in following format:
```
intent_name1, intent utterance 1
intent_name1, intent utterance 2
intent_name2, intent utterance 3
intent_name1, intent utterance 4
...
```
