export const MAX_INT64 = BigInt('9223372036854775807');

export const generateIdempotencyKey = (): string => {
  const chars = 'abcdefghijklmnopqrstuvwxyz0123456789';
  let result = '';
  for (let i = 0; i < 16; i++) {
    result += chars.charAt(Math.floor(Math.random() * chars.length));
  }
  return result;
};

export const clampToInt64 = (digits: string): string => {
  if (!digits) return '';

  let value = BigInt(digits);

  if (value > MAX_INT64) {
    value = MAX_INT64;
  }

  return value.toString();
};

export const formatFromDigits = (digits: string): string => {
  if (!digits) return '';

  const number = Number(digits) / 100;

  return number.toLocaleString('en-US', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  });
};

export const parseMoneyToCents = (value: string): number => {
  return Number(value.replace(/\D/g, ''));
};

export const formatWithCursor = (
  rawValue: string,
  cursorPos: number
): { value: string; cursor: number } => {
  const digitsOnly = rawValue.replace(/\D/g, '');

  const digitsBeforeCursor = rawValue
    .slice(0, cursorPos)
    .replace(/\D/g, '').length;

  const clampedDigits = clampToInt64(digitsOnly);
  const formatted = formatFromDigits(clampedDigits);

  let digitCount = 0;
  let newCursor = 0;

  for (let i = 0; i < formatted.length; i++) {
    if (/\d/.test(formatted[i])) {
      digitCount++;
    }

    if (digitCount >= digitsBeforeCursor) {
      newCursor = i + 1;
      break;
    }
  }

  return {
    value: formatted,
    cursor: newCursor || formatted.length,
  };
};
