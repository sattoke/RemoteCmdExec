# RemoteCmdExec
RemoteCmdExecサーバ上のローカルコマンドをWebサーバを通じてリモートから実行するための簡易的なソフトウェアである。

# 初期設定
makeを実行すると (WindowsならWSLを想定) で実行ファイル (WindowsならRemoteCmdExec.exe) が作られる。
加えて実行ファイルと同じディレクトリに下記のような内容のconfig.yamlという名前のファイルを置く。

```yaml
web:
  address: "0.0.0.0"
  port: 8888
  useTLS: true

tls:
  certFile: "/path/to/cert.pem"
  keyFile: "/path/to/private.pem"

commands:
  - name: "Print Date"
    command: "date"
    params: []

  - name: "Custom Command"
    command: "echo"
    params:
      - "Hello, World!"
```

# 使用方法
RemoteCmdExecの実行ファイル(WindowsならRemoteCmdExec.exe)を実行し、Webブラウザで `config.yaml` で指定したアドレス (https://192.168.1.1:8888 など) にアクセスすると、コマンドのリストが表示されるのでコマンドをクリックすると引数と共にそのコマンドが実行される。
