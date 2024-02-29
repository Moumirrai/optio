export function formatDate(dateStr: string) {
  const date = new Date(dateStr);
  return `${date.toLocaleDateString()} ${date.toLocaleTimeString([], {
    hour: "2-digit",
    minute: "2-digit",
  })}`;
}

export function formatSize(size: number) {
  const i = Math.floor(Math.log(size) / Math.log(1024));
  const sizes = ["B", "KB", "MB", "GB", "TB"];
  return Number((size / Math.pow(1024, i)).toFixed(2)) + " " + sizes[i];
}

export function formatTime(time: number) {
  const msInSecond = 1000;
  const msInMinute = msInSecond * 60;
  const msInHour = msInMinute * 60;

  const hours = Math.floor(time / msInHour);
  time %= msInHour;

  const minutes = Math.floor(time / msInMinute);
  time %= msInMinute;

  const seconds = Math.floor(time / msInSecond);
  time %= msInSecond;

  const milliseconds = time;

  let result = '';
  if (hours > 0) result += `${hours} h `;
  if (minutes > 0) result += `${minutes} m `;
  if (seconds > 0) result += `${seconds} s `;
  if (milliseconds > 0) result += `${milliseconds} ms`;

  return result.trim();
}

export function formatFramerate(input: string): string {
  if (input.includes('/')) {
    // If "/" is present, split the input into two numbers and perform division
    const [num1, num2] = input.split('/');
    const result = parseInt(num1, 10) / parseInt(num2, 10);
    return result.toFixed(2); // Round to 2 decimals
  } else {
    // If "/" is not present, just return the number rounded to 2 decimals
    return parseFloat(input).toFixed(2);
  }
}
