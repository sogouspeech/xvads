# xvads

## 相关说明

eXtended Voice Activity Detection Splitter

非常简单的音频切片调用小工具，用于知音平台相关演示目的。
调用方法类似 xargs，不同的是 xargs 切割文本；xvads 切割音频

- 使用 WebRTC 的 VAD 模块做音频切片
  - WebRTC VAD 的C代码复制自 [py-webrtcvad](https://github.com/wiseman/py-webrtcvad)
  - Go 部分代码和相关封装方法来自 [go-webrtcvad](https://github.com/maxhawkins/go-webrtcvad)

## 构建方法

- 安装 golang 1.13 以上版本
- 较新版本的 gcc 


go install ./cmd/...


## 使用方法

下面是一个示例：

- 用 [fmedia](https://stsaz.github.io/fmedia) 工具抓取默认输入设备的音轨（单声道 16KHZ 16Bit 采样率）并通过管道发送给 xvads
- xvads 将有效音频切割出来并实时调用子命令，将分片音频作为子命令的标准输入

```
fmedia --record --out=@stdout.wav --format=int16 --channels=mono --rate=16000 2>/dev/null | xvads sub-command arguments-of-subcommand ...
```
