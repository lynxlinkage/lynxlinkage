const DATE_FMT = new Intl.DateTimeFormat('en-US', {
	year: 'numeric',
	month: 'short',
	day: 'numeric'
});

export function formatDate(iso: string): string {
	const d = new Date(iso);
	if (Number.isNaN(d.getTime())) return iso;
	return DATE_FMT.format(d);
}

export function employmentLabel(value: string): string {
	switch (value) {
		case 'full_time':
			return 'Full-time';
		case 'part_time':
			return 'Part-time';
		case 'contract':
			return 'Contract';
		case 'internship':
			return 'Internship';
		default:
			return value;
	}
}

export function tierLabel(value: string): string {
	switch (value) {
		case 'strategic':
			return 'Strategic partners';
		case 'exchange':
			return 'Exchanges';
		case 'broker':
			return 'Brokers';
		case 'tech':
			return 'Technology partners';
		default:
			return value;
	}
}

export function sourceLabel(value: string): string {
	switch (value) {
		case 'medium':
			return 'Read on Medium';
		case 'internal':
			return 'Read the article';
		case 'external':
			return 'Read the article';
		default:
			return 'Read';
	}
}
