### Example project for Workstation

#### How to use?

- `$ go run main.go` 

##### Output

```shell
2021-01-17 17:19:05.9828747 +0300 MSK m=+2.345115501 3 3 go,process,19666,720,second_process
2021-01-17 17:19:05.9890697 +0300 MSK m=+2.351310401 1 6 go,process,19666,720,first_process
2021-01-17 17:19:07.9838974 +0300 MSK m=+4.346138101 3 3 go,process,19666,720,second_process
2021-01-17 17:19:07.9895903 +0300 MSK m=+4.351831101 1 6 go,process,19666,720,first_process
2021-01-17 17:19:09.9843827 +0300 MSK m=+6.346623501 3 3 go,process,19666,720,second_process
2021-01-17 17:19:09.9897926 +0300 MSK m=+6.352033401 1 6 go,process,19666,720,first_process
2021-01-17 17:19:11.9850753 +0300 MSK m=+8.347316001 3 3 go,process,19666,720,second_process
2021-01-17 17:19:11.9900587 +0300 MSK m=+8.352299501 1 6 go,process,19666,720,first_process
```

- open another terminal
- `$ ps ax | grep process`

```shell
$ ps ax | grep process
  841 tty1     Sl     0:00 go run ./sub/process.go -id second_process
  843 tty1     Sl     0:00 go run ./sub/process.go -id first_process
 1049 tty1     Sl     0:00 /tmp/go-build536323427/b001/exe/process -id second_process
 1055 tty1     Sl     0:00 /tmp/go-build737399017/b001/exe/process -id first_process
 1064 tty2     S      0:00 grep --color=auto process
```

- `$ kill <pid>` - kill one of sub processes, example `$ kill 1055`
- look in main terminal:

```shell
2021-01-17 17:21:56.8954239 +0300 MSK m=+9.771778301 1 6 Good Luck
2021/01/17 17:21:56 [ERROR]  The async task #first_process process PID 1469 was terminated with an error, the task is removed from the pool and will be restarted in the future. Previous error '%!s(<nil>)'
2021/01/17 17:21:56 [INFO] Async task #first_process sub process PID 1469 successfully killed
2021/01/17 17:21:56 [ERROR]  Reading from stdout for asynctask #first_process completed (with error), cause EOF
```
