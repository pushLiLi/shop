export function useNotificationSound() {
  let audioContext = null

  const play = () => {
    try {
      if (!audioContext) {
        audioContext = new (window.AudioContext || window.webkitAudioContext)()
      }
      if (audioContext.state === 'suspended') {
        audioContext.resume()
      }

      const playTone = (freq, startTime, duration) => {
        const oscillator = audioContext.createOscillator()
        const gainNode = audioContext.createGain()
        oscillator.connect(gainNode)
        gainNode.connect(audioContext.destination)
        oscillator.type = 'sine'
        oscillator.frequency.setValueAtTime(freq, startTime)
        gainNode.gain.setValueAtTime(0.15, startTime)
        gainNode.gain.exponentialRampToValueAtTime(0.01, startTime + duration)
        oscillator.start(startTime)
        oscillator.stop(startTime + duration)
      }

      const now = audioContext.currentTime
      playTone(880, now, 0.15)
      playTone(1100, now + 0.12, 0.2)
    } catch (e) {
      // Audio not available
    }
  }

  return { play }
}
