module.exports = {
    apps: [
        {
            name: "myapp",        // Nama aplikasi di PM2
            script: "./myapp",    // File yang akan dijalankan
            instances: 1,         // Jumlah instance (bisa diubah jadi cluster mode)
            exec_mode: "fork",    // Mode eksekusi (fork untuk aplikasi biasa)
            watch: false,         // Nonaktifkan watch (bisa true untuk auto-restart saat file berubah)
            autorestart: true     // Auto-restart jika aplikasi crash
        }
    ]
};