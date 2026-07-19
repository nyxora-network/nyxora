# مشارکت در NYXORA

اول از همه، ممنون که به مشارکت در NYXORA فکر می‌کنید! ما مشارکت همه را خوشآمد می‌گوییم.

## قوانین رفتاری

این پروژه از [قوانین رفتاری](CODE_OF_CONDUCT.md) پیروی می‌کند. با شرکت، انتظار می‌رود این قوانین را رعایت کنید.

## چگونه می‌توانم مشارکت کنم؟

### گزارش باگ

قبل از ارسال گزارش باگ:
- بررسی کنید آیا قبلاً در [issues](https://github.com/nyxora-network/nyxora/issues) گزارش شده
- اطلاعات جمع‌آوری کنید: نسخه سیستم‌عامل، نسخه Go، مراحل بازتولید، خروجی خطا

**گزارش باگ ارسال کنید** با باز کردن [issue جدید](https://github.com/nyxora-network/nyxora/issues/new?template=bug_report.md).

### پیشنهاد ویژگی جدید

[درخواست ویژگی](https://github.com/nyxora-network/nyxora/issues/new?template=feature_request.md) باز کنید و توضیح دهید:
- مشکلی که حل می‌کنید
- راه‌حل مورد نظر شما
- جایگزین‌های در نظر گرفته شده

### اضافه کردن ترنسپورت جدید

1. فایل `internal/transport/<name>.go` ایجاد کنید و interface `Transport` را پیاده‌سازی کنید
2. آن را در `internal/transport/registry.go` ثبت کنید
3. پوشه `tunnels/<name>/` با اسکریپت‌های نصب ایجاد کنید
4. وزن امتیازدهی در `internal/transport/scoring.go` اضافه کنید
5. تست بنویسید و `make test` اجرا کنید

### بهبود TUI

TUI تعاملی در `internal/interactive/` قرار دارد و از این کتابخانه‌ها استفاده می‌کند:
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) — فریمورک TUI
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) — استایل‌بندی
- [Bubbles](https://github.com/charmbracelet/bubbles) — کامپوننت‌ها (textinput, spinner, progress)

رنگ‌های تم از مقادیر TrueColor Catppuccin در فایل `theme.go` استفاده می‌کنند.

## فرآیند Pull Request

1. ریپازیتوری را Fork کنید و شاخه خود را از `main` بسازید
2. `make test` و `make vet` اجرا کنید — هر دو باید موفق باشند
3. برای ویژگی‌های جدید تست اضافه کنید
4. در صورت نیاز مستندات را به‌روز کنید
5. مطمئن شوید کد شما از قراردادهای موجود پیروی می‌کند
6. PR را با توضیحات واضح ارسال کنید

### سبک Commit

از پیام‌های commit استاندارد استفاده کنید:
- `feat:` — ویژگی جدید
- `fix:` — رفع باگ
- `refactor:` — تغییر کد بدون رفع باگ/ویژگی جدید
- `docs:` — فقط مستندات
- `test:` — اضافه/رفع تست
- `style:` — قالب‌بندی، تغییرات استایل
- `chore:` — نگهداری، وابستگی‌ها

## راه‌اندازی محیط توسعه

```bash
# Fork و clone
git clone https://github.com/YOUR_USERNAME/nyxora.git
cd nyxora

# اضافه کردن upstream remote
git remote add upstream https://github.com/nyxora-network/nyxora.git

# ایجاد شاخه ویژگی
git checkout -b feat/your-feature

# تغییرات اعمال کنید، سپس:
make test
make vet
make build

# Commit و push
git commit -m "feat: add your feature"
git push origin feat/your-feature
```

## سوالی دارید؟

[بحث](https://github.com/nyxora-network/nyxora/discussions) باز کنید یا به [تلگرام](https://t.me/NyxoraCore) بپیوندید.
