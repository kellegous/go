export function restartableAsync<T>(iter: AsyncIterable<T>): AsyncIterable<T> {
  // buffer stores all items that have been previously consumed.
  const buffer: T[] = [];
  const gen = async function* () {
    // index of the next item in the buffer to yield.
    let i = 0;
    // produce all items previously consumed by other iterators.
    for (; i < buffer.length; i++) {
      yield buffer[i];
    }
    // now takes the next from the iterator.
    for await (const item of iter) {
      // this is a little subtle, but other concurrent iterators may have
      // consumed and buffered items while we were waiting. So we need to put
      // our new item in the back of the buffer and yield from where we previously
      // left off.
      buffer.push(item);
      for (; i < buffer.length; i++) {
        yield buffer[i];
      }
    }
  };

  return {
    [Symbol.asyncIterator]: () => gen(),
  };
}
