/// <reference types="vitest" />

import { render, screen } from '@testing-library/react';
import App from '../src/App';

describe('App', () => {
	it('renders the welcome hero', () => {
		render(<App />);
		expect(
			screen.getByRole('heading', {
				name: /frontend scaffolding ready to go/i,
			}),
		).toBeInTheDocument();
		expect(
			screen.getByText(/Tailwind CSS, Biome, and Vitest are wired up/i),
		).toBeInTheDocument();
	});
});
