export type ResearchSource = 'medium' | 'internal' | 'external';

export interface ResearchCard {
	id: number;
	title: string;
	summary: string;
	tags: string[];
	coverImageUrl?: string;
	externalUrl: string;
	source: ResearchSource;
	publishedAt: string;
	displayOrder: number;
}

export type EmploymentType = 'full_time' | 'part_time' | 'contract' | 'internship';

export interface JobPosting {
	id: number;
	title: string;
	team: string;
	location: string;
	employmentType: EmploymentType;
	descriptionMd: string;
	applyUrlOrEmail: string;
	postedAt: string;
	isActive: boolean;
	createdAt?: string;
	updatedAt?: string;
	createdBy?: number;
	updatedBy?: number;
}

export type PartnerTier = 'strategic' | 'exchange' | 'broker' | 'tech';

export interface Partner {
	id: number;
	name: string;
	logoUrl: string;
	websiteUrl?: string;
	tier: PartnerTier;
	description?: string;
	displayOrder: number;
}

export type ContactKind = 'general' | 'partnership' | 'research' | 'hiring';

export interface ContactPayload {
	name: string;
	email: string;
	company?: string;
	message: string;
	kind?: ContactKind;
}

export interface ListResponse<T> {
	items: T[];
}

export type Role = 'hr';

export interface User {
	id: number;
	email: string;
	role: Role;
	createdAt: string;
	lastLoginAt?: string;
}

export interface JobUpsertPayload {
	title: string;
	team: string;
	location: string;
	employmentType: EmploymentType;
	descriptionMd: string;
	applyUrlOrEmail: string;
	postedAt?: string;
	isActive?: boolean;
}
