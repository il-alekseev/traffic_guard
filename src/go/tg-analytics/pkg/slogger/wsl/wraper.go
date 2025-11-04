package wsl

import "log/slog"

func Label(labelName string, value string) slog.Attr {
	return slog.Any(labelName, value)
}

func Int(valueName string, intValue int) slog.Attr {
	return slog.Int(valueName, intValue)
}

func Int64(valueName string, int64Value int64) slog.Attr {
	return slog.Int64(valueName, int64Value)
}

func String(valueName string, strValue string) slog.Attr {
	return slog.String(valueName, strValue)
}

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}

func Info(info string) slog.Attr {
	return slog.String("msg", info)
}
