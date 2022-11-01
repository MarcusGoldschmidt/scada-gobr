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