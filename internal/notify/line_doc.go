// Package notify provides notification integrations for cronwrap.
//
// LINE Notify
//
// LineNotifier sends job result notifications via the LINE Notify API
// (https://notify-bot.line.me/). A personal access token obtained from
// the LINE Notify service is required.
//
// Usage:
//
//	n := notify.NewLineNotifier(os.Getenv("LINE_NOTIFY_TOKEN"))
//	n.Notify("my-cron-job", err)
package notify
