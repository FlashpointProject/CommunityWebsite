export function easyDateTimeFormat(date: Date | string): string {
  if (typeof date === 'string') {
    date = new Date(date);
  }
  const options: Intl.DateTimeFormatOptions = {
    year: '2-digit',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    hour12: true
  };

  return new Intl.DateTimeFormat('en-US', options).format(date);
}

export function easyDateFormat(date: Date | string): string {
  if (typeof date === 'string') {
    date = new Date(date);
  }
  const options: Intl.DateTimeFormatOptions = {
    year: '2-digit',
    month: '2-digit',
    day: '2-digit',
  };

  return new Intl.DateTimeFormat('en-US', options).format(date);
}
