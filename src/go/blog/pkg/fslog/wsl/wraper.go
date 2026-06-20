package wsl

import "log/slog"

func Label(labelName string, value string) slog.Attr {
	return slog.Any(labelName, value)
}

func Int(valueName string, intValue int) slog.Attr {
	return slog.Int(valueName, intValue)
}

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
