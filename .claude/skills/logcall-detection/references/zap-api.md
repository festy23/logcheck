# go.uber.org/zap API Reference

## *zap.Logger methods

| Method                    | Msg Index | KV Style                        |
|---------------------------|-----------|----------------------------------|
| Debug(msg, fields...)     | 0         | zap.Field (zap.String, zap.Int)  |
| Info(msg, fields...)      | 0         | zap.Field                        |
| Warn(msg, fields...)      | 0         | zap.Field                        |
| Error(msg, fields...)     | 0         | zap.Field                        |
| DPanic(msg, fields...)    | 0         | zap.Field                        |
| Fatal(msg, fields...)     | 0         | zap.Field                        |
| Panic(msg, fields...)     | 0         | zap.Field                        |

KV extraction for Logger: each field is a call like `zap.String("key", val)`.
The key is the first argument of the zap helper function.

## *zap.SugaredLogger methods

| Method                       | Msg Index | KV Style              | IsFormat |
|------------------------------|-----------|------------------------|----------|
| Debug(args...)               | 0         | none (args concat)     | false    |
| Info(args...)                | 0         | none                   | false    |
| Warn(args...)                | 0         | none                   | false    |
| Error(args...)               | 0         | none                   | false    |
| DPanic(args...)              | 0         | none                   | false    |
| Fatal(args...)               | 0         | none                   | false    |
| Panic(args...)               | 0         | none                   | false    |
| Debugf(template, args...)    | 0         | none                   | true     |
| Infof(template, args...)     | 0         | none                   | true     |
| Warnf(template, args...)     | 0         | none                   | true     |
| Errorf(template, args...)    | 0         | none                   | true     |
| DPanicf(template, args...)   | 0         | none                   | true     |
| Fatalf(template, args...)    | 0         | none                   | true     |
| Panicf(template, args...)    | 0         | none                   | true     |
| Debugw(msg, keysAndVals...)  | 0         | alternating string/any | false    |
| Infow(msg, keysAndVals...)   | 0         | alternating string/any | false    |
| Warnw(msg, keysAndVals...)   | 0         | alternating string/any | false    |
| Errorw(msg, keysAndVals...)  | 0         | alternating string/any | false    |
| DPanicw(msg, keysAndVals...) | 0         | alternating string/any | false    |
| Fatalw(msg, keysAndVals...)  | 0         | alternating string/any | false    |
| Panicw(msg, keysAndVals...)  | 0         | alternating string/any | false    |

## Methods to IGNORE (no message)

With(), Named(), Sugar(), Desugar(), Sync(), Core(), Check(), WithOptions()
