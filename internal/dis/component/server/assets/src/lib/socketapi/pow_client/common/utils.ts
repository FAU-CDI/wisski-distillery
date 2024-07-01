/** resolves the promise after the specified amount of milliseconds */
export async function sleep (milliseconds: number): Promise<void> {
  await new Promise<void>(resolve => setTimeout(resolve, milliseconds))
}
