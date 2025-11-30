I'm having this idea on a whim, so it's probably bad.

Imagine a color wheel, with some points arbitrarily assigned to "ideal" waveforms. Different sounds are often described as having different colors, so one such mapping would describe this relationship for maybe one genre of music. 

Waveforms can be approximated using a harmonic series. For an arbitrary input, the harmonic series extracted via a fourier transform can be aligned, as best possible, with each of the ideal waveforms, and the "contribution" of the associated color is the quality or error of that match.

The matching is done by changing the pitch and amplitude of the harmonic series of the ideal form, and comparing that to the fourier transform of the input. The quality of the match is a factor of the similarity and amplitude of the matched waveform. The color on the color wheel can be visualized as glowing proportionately to how well its matching to the input at any given time.

The goal is to 1. produce a light show- a nice color-wheel visualization of the sound. And 2. Give the user a target that they can reach how they please. The pitch of the ideal frequencies should, ignoring interference patterns, disappear during matching, only measuring the tambre or "color" of the sound as an error against the reference pattern. The user is pointed to parts of the circle and asked to adjust their circuit, as well as route inputs, to match the color.
# Questions
First, does anything like this already exist?
 - A quick search of google and google scholar says no. But that's not conclusive.
 - ChatGPT has a different story to tell. Specifically, I was able to ask about "basis-dependent" timbre descriptors, and that came up with a few science-y looking things. Might be worth more research later. https://chatgpt.com/share/69269718-4150-8003-b357-638858c164de
	 - The first two examples are ok, but they miss the point I'm trying to make with the arbitrary colors.
		 - MFCCs can sometimes do a better job at encoding frequencies in a way that preserves color, and so they kind of encode the color itself.
		 - Spectral centroid measures brightness on a spectrum.
	 - Then some really interesting stuff
		 - Low dimensional timbre spaces. Grey & Gordon 1978. This is pretty much the same thing I was talking about, except real! Never know how much it could actually help, but good to have on hand. Here's a rad chart https://www.at.or.at/hans/misc/timbre-space/timbre-space.html
		 - STFT. Family of extensions to the basic FFT idea with special functions in describing the time window, amplitude range, and other aspects of the classic fourier transform. Can be used to reveal more meaningful shapes.
		 - Wavelet basis. Best modern example from the looks of things. Choice of wave as a basis isn't completely arbitrary, I'd have to do a lot of reading to see what's allowed or makes sense here. Also parameterized by a Q-factor octave resolution and other tunings.
Second, is this worth your time?
 - No, not yet. In the long term, I wanted a way to create more open-ended puzzles to make the project more appealing. However, that's only one part of appeal that puzzle games struggle with. The other is with being weird, and I'd be doing myself no favors there.