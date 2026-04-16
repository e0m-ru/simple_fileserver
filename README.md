# 📡 Local File Sharer

A simple utility for sharing files on a local network.  
Automatically prompts you to select a network interface and generates a QR code for quick connection.

## ✨ Opportunities
- 🌐 Select local network interface
- 📱 QR code for instant connection from your phone or other PC
- ⚡ Works without external dependencies
- 🔒 Files are transferred only within your network

## 🛠 Usage

```bash
./fileserver [options]
```

| Flag | Description | Default |
| ---- | -- | -- |
| `-p` | File distribution port | `8080` |
| `-u` | Folder for saving downloaded files | OS temporary directory |

## 💡 Examples

```bash
# Run with default settings
./filesrever

# Run on port 3000, files are saved to ~/Downloads
./filesrever -p 3000 -u ~/Downloads
```

## 📦 Assembly (if necessary)

```bash
go build -o filesrever .
```

## 📝 Notes
- To stop the server, press `Ctrl+C`.
- Make sure the firewall allows incoming connections on the selected port.
```