# 🎮 Perfect World - Discord Rich Presence (RPC) & Game Starter

Standalone Game Starter & Discord Rich Presence (RPC) Client untuk server **Perfect World: Realm of Chaos (v1.4.6 Sirens of War)**.

Aplikasi ini berfungsi sebagai launcher instan yang otomatis menghubungkan profil Discord pemain dengan game Perfect World saat mereka bermain, menampilkan status interaktif, logo server, durasi bermain, serta tombol link ke website dan Discord resmi.

---

## ✨ Fitur Utama

- 🚀 **Zero Dependencies**: Berjalan langsung tanpa perlu menginstal runtime tambahan (.NET Framework / Python runtime).
- 💬 **Discord Rich Presence Otomatis**:
  - Menampilkan status *"Playing Realm of Chaos"*.
  - Menampilkan logo server dan versi client (`v1.4.6 build 2305`).
  - Menghitung durasi bermain secara real-time (*Elapsed Time*).
  - Tombol interaktif klik langsung: **🌐 Website** (`https://your-server-website.com`) & **💬 Discord Server**.
- ⚙️ **Fully Configurable (`config.json`)**: Seluruh teks status, gambar, tombol link, dan parameter eksekusi game dapat diatur tanpa perlu compile ulang.
- 🛡️ **Safe & Native**: Tidak melakukan injeksi memory/hooking ilegal yang memicu alarm palsu antivirus (*false positive*).
- 🔕 **Silent Background**: Berjalan di latar belakang tanpa jendela command prompt hitam yang mengganggu.
- ⚠️ **Native Error Alert**: Memberitahu pemain dengan dialog popup Windows jika diletakkan di luar folder game.

---

## 📁 Struktur Penempatan di Client Game

Letakkan file `RealmOfChaos.exe` dan `config.json` di **folder utama client game** (sejajar dengan folder `element/`):

```text
📁 Perfect World - Realm of Chaos/
├── 📄 RealmOfChaos.exe      <--- File starter ini
├── 📄 config.json           <--- Konfigurasi teks & link
├── 📁 element/
│   ├── 📄 elementclient.exe
│   ├── 📁 data/
│   └── ...
├── 📁 launcher/ (opsional)
└── 📁 patcher/  (opsional)
```

---

## ⚙️ Konfigurasi (`config.json`)

Anda dapat menyesuaikan isi tampilan Discord Rich Presence melalui file `config.json`:

```json
{
  "client_id": "YOUR_DISCORD_APPLICATION_ID",
  "details": "Perfect World v1.4.6",
  "state": "Playing on Realm of Chaos",
  "large_image": "logo_roc",
  "large_text": "Realm of Chaos - Sirens of War",
  "small_image": "pwi",
  "small_text": "v1.4.6 build 2305",
  "buttons": [
    {
      "label": "🌐 Website",
      "url": "https://your-server-website.com"
    },
    {
      "label": "💬 Discord Server",
      "url": "https://discord.gg/your-discord"
    }
  ],
  "game_executable": "element/elementclient.exe",
  "game_arguments": [
    "game:cpw",
    "console:1"
  ],
  "update_interval_seconds": 15
}
```

### Keterangan Pengaturan:
| Parameter | Tipe | Keterangan |
| :--- | :--- | :--- |
| `client_id` | String | Application ID dari Discord Developer Portal. |
| `details` | String | Baris pertama status di profil Discord pemain. |
| `state` | String | Baris kedua status di profil Discord pemain. |
| `large_image` | String | Nama asset gambar logo utama di Discord Developer Portal. |
| `large_text` | String | Teks tooltip saat kursor diarahkan ke logo utama. |
| `small_image` | String | Nama asset gambar badge kecil di sudut logo utama. |
| `small_text` | String | Teks tooltip gambar badge kecil. |
| `buttons` | Array | Tombol link klik langsung (maksimal 2 tombol). |
| `game_executable` | String | Lokasi file binary game (`element/elementclient.exe`). |
| `game_arguments` | Array | Argumen pembuka game (`game:cpw`, `console:1`, dll). |

---

## 🛠️ Kompilasi Ulang (Build from Source)

Jika ingin mengompilasi ulang binary Windows dari Linux / macOS / Windows:

### 1. Prasyarat
- Go 1.18 atau lebih baru.

### 2. Kompilasi untuk Windows (32-bit & 64-bit)

```bash
# Windows 32-bit (kompatibel untuk semua versi Windows):
GOOS=windows GOARCH=386 go build -ldflags="-s -w -H windowsgui" -o bin/RealmOfChaos.exe .

# Windows 64-bit:
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w -H windowsgui" -o bin/RealmOfChaos_x64.exe .
```

*Catatan Flag `-ldflags="-s -w -H windowsgui"`:*
- `-s -w`: Menghapus simbol debug untuk memperkecil ukuran file secara maksimal.
- `-H windowsgui`: Menghilangkan tampilan hitam konsol/CMD saat program dijalankan.

---

## 📝 Lisensi
Dilisensikan untuk server **Perfect World Realm of Chaos** © 2026.
