<?php
// シナリオデータの定義（ID、テキスト、背景、選択肢、フラグ変化）
$scenarioMaster = [
    "start" => [
        "text" => "目の前に古い格納庫がある。静寂が満ちている。",
        "bg" => "#1a1a2e",
        "choices" => [
            ["text" => "扉を開けて中に入る", "next" => "inside", "set_flag" => "entered"],
            ["text" => "その場を立ち去る", "next" => "leave", "set_flag" => null]
        ]
    ],
    "inside" => [
        "text" => "中には退役した戦闘機が眠っていた。計器類はまだ生きているようだ。",
        "bg" => "#16213e",
        "choices" => [
            ["text" => "コックピットに乗り込む", "next" => "cockpit", "set_flag" => null]
        ]
    ],
    "leave" => [
        "text" => "君は背を向けた。ミッションはここで終了だ。（BAD END）",
        "bg" => "#0f0f16",
        "choices" => []
    ],
    "cockpit" => [
        "text" => "シートに深く腰掛ける。センサーフュージョンが起動し、視界に情報が同期された。（TRUE END）",
        "bg" => "#0f3460",
        "choices" => []
    ]
];

// Go言語やJSが最速でパースできるようにJSONとして書き出す
$outputPath = __DIR__ . 'scenario.json';
file_put_contents($outputPath, json_encode($scenarioMaster, JSON_UNESCAPED_UNICODE | JSON_PRETTY_PRINT));

echo "[PHP] シナリオマスタ(scenario.json)を正常に出力しました。\n";