import { Mutex, MutexInterface } from 'async-mutex'

const error = console.error.bind(console)

/** discard discards the result of a promise, or logs an error if it occurs */
export function discard<T> (p: Promise<T>): void {
  p.then((): void => {}).catch(error)
}

/** runs worker exclusively on m, and discards the resulting promise */
export function runMutexExclusive<T> (m: Mutex, worker: MutexInterface.Worker<T>): void {
  discard(m.runExclusive(worker))
}
