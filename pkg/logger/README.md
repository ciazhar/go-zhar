# Logging

Paling bener ngelog ke stdout emang. Mau diforward ke file atau diingest oleh log collector gampang

./binary >> logfile.log 2>&1

dia klaim "bikin log segampang ini loh, ga perlu depend ke community module/package"
padahal implementasi dia banyak flawnya. At least concern dia valid sih, bahwa logger memang harus concurrent safe (bisa diakses multiple goroutine secara bersamaan). Cuma apakah harus pake mutex, ini yang jadi pertanyaan.
Flawnya di mana? Logging ga cuma harus concurrent safe, tapi juga performant. Consider log outputnya ke resource yang harus dimanage IOnya. Ke file aja ga bisa tiap logging langsung write ke file. Biasanya ditampung di buffer dulu, setelah sekian bytes diflush. Logging ke centralized log (remote URL) juga sama. Belum gimana flush ketika graceful shutdown

https://www.facebook.com/groups/GophersID/permalink/7645193708833780/