# airy-wave

Airy（线性/小振幅）波浪理论计算工具。

在 Airy 模型中，均匀水深 h 中周期 T、振幅 a 的推进波满足色散关系：

```
ω² = g · k · tanh(k·h)
```

据此可求得波长 λ = 2π/k、相速 c = λ/T，以及任意位置、任意水深处的水质点水平/垂直速度和动水压强。本项目提供：

- `internal/wave` — 波数（迭代求解色散关系）、波长、相速、水面高程与轨道速度、动水压强。
- `internal/fluid` — 沿垂直剖面采样并导出 CSV。
- `internal/server` — HTTP 接口，提供网页控制台与 JSON 接口。

## 用法

```bash
go run . -http :8080            # 在 :8080 提供网页控制台与 /api
go run . -http :8080 -example-dir example
```

浏览器打开 http://localhost:8080 即可交互：网页载入内置示例（振幅 1 m、周期 8 s、水深 50 m），调用 `/api/wave` 并在垂向剖面绘制水平速度分布。计算全部在本地进程内完成。

### HTTP 接口

- `POST /api/wave` — 计算波要素与垂直剖面。
  请求：`{"amplitude":1,"period":8,"depth":50,"layers":12,"rho":1025}`。
  返回：`{"wavelength":99.02,"wavenumber":0.0634,"phase_speed":12.38,"deep_water":false,"profile":[...]}`。
- `/example/` 提供示例 JSON 文件。

## 结论

对 T=8 s、h=50 m 的波浪，深水判据 h > 0.5λ 近似成立（λ ≈ 99 m），相速约 12.4 m/s；水质点呈椭圆轨道，随深度增加速度幅值按 cosh 衰减、动水压强按 cosh 衰减。

## 构建与测试

```bash
go build ./...
go test ./...
```

## 许可

见 LICENSE。
