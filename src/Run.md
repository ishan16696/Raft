## Demo

### `$go build main.go`

#### `$./main --help`
```
Usage of ./main:
  -nodes int
    	total no. of nodes (default 3)
  -peers string
    	specify the portNo. of rest of all peer nodes (default "9091,9092")
  -port int
    	specify the portNo. of node (default 9090)
```

### start the 3 terminal with 

`$./main -nodes=3 -peers="9091,9092" -port=9090`
`$./main -nodes=3 -peers="9092,9090" -port=9091`
`$./main -nodes=3 -peers="9091,9090" -port=9092`

