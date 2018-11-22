# updatedb-go
为IPIP.net付费客户提供的下载与更新工具，支持ipdb与文本格式。

# Linux
<pre>
git clone https://github.com/ipipdotnet/updatedb-go
cd updatedb-go/cmd
ipipdowner-linux-x64
</pre>

# Windows

<pre>
git clone https://github.com/ipipdotnet/updatedb-go
cd updatedb-go/cmd
ipipdowner-windows-x64.exe
</pre>

# Help
<pre>
Example:
        下载ipdb格式数据文件
        ipipdowner-windows-x64.exe --dir=c:/temp --type=ipdb --token=XXX
        下载文本格式数据文件
        ipipdowner-windows-x64.exe --dir=c:/temp --type=txtx --token=XXX

Usage of ipipdowner-windows-x64.exe:
      --compress       --compress (default true)
      --dir string     --dir=/tmp
      --lang string    --lang=EN|CN
      --merge          --merge
      --token string   --token=XXX
      --type string    --type=ipdb|txtx (default "ipdb")
      --view           --view show download url
</pre>      