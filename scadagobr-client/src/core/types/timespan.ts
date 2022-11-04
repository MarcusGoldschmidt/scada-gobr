export function validateTimespan(value: string): boolean {
    try {
        const [h, m, s] = value.split(":")

        return !(+h > 23 || +m > 59 || +s > 59);
    } catch (e) {
        return false;
    }
}

export function validateDuration(value: string): boolean {
    try {
        const h = +(value[0] + value[1])
        const m = +(value[3] + value[4])
        const s = +(value[6] + value[7])

        return !(+h > 23 || +m > 59 || +s > 59);
    } catch (e) {
        return false;
    }
}

export function timespanToNanosecond(value: string): number {
    const [h, m, s] = value.split(":")

    return ((+h * 60 * 60) + (60 * +m) + +s) * 1000000000
}

export function getTimeFormNanosecond(value: number): [number, number, number] {
    const seconds = value / 1_000_000_000

    const h = Math.floor(seconds / 3600)
    const m = Math.floor((seconds % 3600) / 60)
    const s = Math.floor(seconds % 60)

    return [h, m, s]
}

export function nanosecondToTimespan(value?: number): string {
    if (value === undefined) {
        return ""
    }

    const [h, m, s] = getTimeFormNanosecond(value)

    return `${h.toString().padStart(2, "0")}:${m.toString().padStart(2, "0")}:${s.toString().padStart(2, "0")}`
}

export function nanosecondToTimespanFormatted(value?: number, fallback: string = "00h00m00s"): string {
    if (value === undefined) {
        return fallback
    }

    const [h, m, s] = getTimeFormNanosecond(value)

    return `${h.toString().padStart(2, "0")}h:${m.toString().padStart(2, "0")}m:${s.toString().padStart(2, "0")}s`
}