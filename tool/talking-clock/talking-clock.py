#!/usr/bin/env python3
"""
根据当前时间生成报时文本，使用 edge-tts 合成语音后调用 su-home audio API 播放。

示例:
  pip install -r test/audio/requirements.txt
  python tool/talking-clock/talking-clock.py
  SU_HOME_DEVICE=linux SU_HOME_API=http://127.0.0.1:41406 python tool/talking-clock/talking-clock.py
"""

from __future__ import annotations

import argparse
import asyncio
from datetime import datetime
import os
import tempfile

import edge_tts
import requests

WEEKDAYS = ("星期一", "星期二", "星期三", "星期四", "星期五", "星期六", "星期日")


def build_time_text(now: datetime) -> str:
    period = "凌晨"
    if 5 <= now.hour < 12:
        period = "上午"
    elif 12 <= now.hour < 14:
        period = "中午"
    elif 14 <= now.hour < 18:
        period = "下午"
    elif 18 <= now.hour < 24:
        period = "晚上"
    hour = now.hour
    if hour > 12:
        hour -= 12
    return (
        f"现在是"
        # f"{now.year}年{now.month}月{now.day}日，"
        # f"{WEEKDAYS[now.weekday()]}，"
        f"{period}{hour}点{now.minute}分。"
    )


async def synthesize_to_file(text: str, voice: str, path: str) -> str:
    communicate = edge_tts.Communicate(text, voice)
    await communicate.save(path)
    return path


def play_on_device(api_base: str, device: str, audio_path: str) -> dict:
    url = f"{api_base.rstrip('/')}/api/devices/{device}/audio/play"
    name = os.path.basename(audio_path)
    with open(audio_path, "rb") as f:
        resp = requests.post(
            url,
            files={"audiofile": (name, f, "audio/mpeg")},
            timeout=300,
        )
    resp.raise_for_status()
    try:
        return resp.json()
    except ValueError:
        return {"raw": resp.text}


async def run(text: str, voice: str, api_base: str, device: str) -> None:
    fd, path = tempfile.mkstemp(prefix="talking-clock-", suffix=".mp3")
    os.close(fd)
    try:
        print(f"Text: {text}")
        print(f"TTS: voice={voice!r}, chars={len(text)}")
        await synthesize_to_file(text, voice, path)
        print(f"Saved: {path} ({os.path.getsize(path)} bytes)")
        print(f"API: POST .../devices/{device}/audio/play")
        body = play_on_device(api_base, device, path)
        print("Response:", body)
    finally:
        if os.path.isfile(path):
            os.remove(path)


def main() -> int:
    parser = argparse.ArgumentParser(description="当前时间报时 edge-tts -> su-home audio play API")
    parser.add_argument(
        "--text",
        help="覆盖默认报时文本；不传时根据当前本地时间生成",
    )
    args = parser.parse_args()

    api = os.environ.get("SU_HOME_API", "http://127.0.0.1:41406")
    device = os.environ.get("SU_HOME_DEVICE", "linux")
    voice = os.environ.get("EDGE_TTS_VOICE", "zh-CN-XiaoxiaoNeural")
    text = args.text or build_time_text(datetime.now())
    asyncio.run(run(text, voice, api, device))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
