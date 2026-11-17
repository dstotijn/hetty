# Panduan Operator (untuk kamu, bukan untuk agent)

Dokumen di repo ini ditulis dalam bahasa Inggris karena semua model, baik di Claude Code maupun OpenCode, paling
konsisten mengikuti instruksi teknis berbahasa Inggris. Panduan ini menjelaskan cara kamu menjalankan proyeknya.

## 1. Cara pasang paket dokumen ini

Repo ini adalah **fork GitHub dari `dstotijn/hetty`** (sengaja, untuk menghormati pembuat Hetty). Konsekuensinya, `gh`
secara default mengarahkan PR ke repo Hetty. Kunci dulu supaya selalu ke repo-mu:

```bash
# Fork di GitHub (centang "Copy the main branch only"), rename repo jadi "ratatoskr", lalu:
git clone https://github.com/USERNAME/ratatoskr.git && cd ratatoskr
git remote add upstream https://github.com/dstotijn/hetty.git
git remote set-url --push upstream DISABLED      # mencegah push ke Hetty
gh repo set-default USERNAME/ratatoskr           # gh selalu ke repo-mu
git tag hetty-v0.7.0-base && git push origin --tags
unzip ~/Downloads/ratatoskr-agent-docs.zip -d .
sed -i 's#OWNER/ratatoskr#USERNAME/ratatoskr#g' AGENTS.md docs/*.md docs/tasks/*.md
git add . && git commit -m "docs: add Ratatoskr specs and agent setup" && git push origin main
```
Di GitHub: Settings → Branches → aturan proteksi `main` (wajib PR + CI hijau). Settings → Actions → aktifkan workflow.

Ganti `OWNER` di `docs/tasks/T-001-*.md` (dan nanti `go.mod`) dengan username GitHub-mu. Codename "Ratatoskr" bisa
diganti. Cukup cari dan ganti di `docs/`, `AGENTS.md`, dan file agent.

## 2. Siapa mengerjakan apa

| Peran | Tool yang disarankan | Model |
|---|---|---|
| lead (arsitek, kontrak, perencanaan) | Claude Code | model terkuat (Opus) |
| backend (Go) | Claude Code atau OpenCode | model coding terkuat yang kamu punya |
| frontend (Svelte) | OpenCode atau Claude Code | model coding kuat, boleh beda dari backend |
| reviewer | **tool/model yang berbeda dari implementer** | model reasoning kuat |
| qa (menulis test) | OpenCode | model yang lebih murah dan cepat sudah cukup |

Model per agent bisa diatur di frontmatter:
- Claude Code: `.claude/agents/<role>.md` → `model: opus` / `sonnet` / `inherit`
- OpenCode: `.opencode/agents/<role>.md` → `model: <provider>/<model-id>` (lihat daftar dengan `opencode models`)

## 3. Alur kerja harian

1. **Lead merencanakan.** Di Claude Code (root repo, branch `main`), jalankan `/plan-milestone M1` saat milestone
   sebelumnya sudah sekitar 70% selesai. File task untuk M0 dan M1 sudah disiapkan.
2. **Pilih task yang siap.** Task siap dikerjakan kalau semua `depends_on`-nya sudah `done`. Urutan awal:
   T-001, lalu T-002 dan T-003 paralel, lalu T-004, T-005, T-008 paralel, dan seterusnya.
3. **Satu task = satu worktree = satu agent.**
   ```bash
   git worktree add ../ratatoskr-T-010 -b task/T-010-rawhttp main
   cd ../ratatoskr-T-010
   claude            # lalu: /task T-010   (pakai subagent "backend")
   # atau
   opencode          # pilih agent "backend", lalu: /task T-010
   ```
   Kamu bisa menjalankan 3 sampai 5 worktree paralel di terminal terpisah (tmux sangat membantu).
4. **Review silang.** Saat task berstatus `in_review`, buka worktree yang sama dengan tool/model **lain**, lalu jalankan
   `/review T-010`. Kalau hasilnya `changes_requested`, kembali ke implementer dengan `/task T-010`.
5. **Merge.** Squash-merge ke `main` sesuai urutan dependensi, ubah status task jadi `done`, hapus worktree:
   `git worktree remove ../ratatoskr-T-010`.
6. **Kalau agent "blocked"**, baca pertanyaannya di file task, jawab dengan memperbarui spec (atau minta lead yang
   memperbarui), lalu set status kembali ke `todo`/`in_progress`.

## 4. Aturan emas yang perlu kamu jaga sendiri

- Jangan biarkan dua agent mengedit area yang sama. Kolom `paths` di file task adalah pagarnya.
- Perubahan skema GraphQL selalu lewat lead, sebagai PR tersendiri, sebelum task backend/frontend fitur itu dimulai.
- Kalau sesi agent sudah panjang dan mulai melantur, jalankan `/handoff`, lalu mulai sesi baru.
- Jangan merge kalau `make check` merah, walaupun agent bilang "hampir selesai".
- Fitur boleh meniru perilaku Caido, tapi kode, aset, warna, ikon, teks, dan nama tidak boleh disalin.

## 5. Ekspektasi realistis

Dengan 3 sampai 5 agent paralel dan review disiplin, M0 sampai M1 (proxy jalan, history, query, scope, sitemap)
kira-kira butuh beberapa minggu kerja operator. Parity penuh (M7) adalah proyek berbulan-bulan. Prioritaskan M0 sampai
M3, karena itu sudah membuat tool ini berguna sehari-hari.

## 6. Kalau ingin mengubah keputusan teknis

Buat ADR baru di `docs/adr/` (salin `0000-template.md`), tandai ADR lama `superseded`, lalu minta lead memperbarui
spec dan task yang terdampak.
