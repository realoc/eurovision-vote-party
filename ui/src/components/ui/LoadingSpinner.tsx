export function LoadingSpinner() {
	return (
		<output
			aria-label="Loading"
			className="flex items-center justify-center p-6 text-indigo-500"
		>
			<svg
				className="h-6 w-6 animate-spin"
				viewBox="0 0 24 24"
				fill="none"
				role="img"
				aria-label="Spinner"
				xmlns="http://www.w3.org/2000/svg"
			>
				<circle
					className="opacity-20"
					cx="12"
					cy="12"
					r="10"
					stroke="currentColor"
					strokeWidth="4"
				/>
				<path
					className="opacity-80"
					d="M22 12c0-5.523-4.477-10-10-10"
					stroke="currentColor"
					strokeWidth="4"
					strokeLinecap="round"
				/>
			</svg>
		</output>
	);
}

export default LoadingSpinner;
