## Demo

1. Compile: `$go build main.go`
2.  `$./main --help`
    ```
    Usage of ./main:
      -nodes int
          total no. of nodes (default 3)
      -peers string
          specify the portNo. of rest of all peer nodes (default "9091,9092")
      -port int
          specify the portNo. of node (default 9090)
    ```

3. Start the 3 terminal: 

`$./main -nodes=3 -peers="9091,9092" -port=9090` <br/>
`$./main -nodes=3 -peers="9092,9090" -port=9091` <br/>
`$./main -nodes=3 -peers="9091,9090" -port=9092` <br/>

<img width="1608" alt="Screenshot 2021-10-24 at 11 29 14 PM" src="https://user-images.githubusercontent.com/42602577/138606722-6ae2b478-c133-4bb6-9ed2-146523f59ca4.png">
