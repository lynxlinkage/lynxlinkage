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

export interface ApplicationFile {
	id: number;
	applicationId: number;
	originalName: string;
	contentType: string;
	sizeBytes: number;
	createdAt: string;
}

export type ApplicationStatusKind = 'open' | 'accept' | 'reject';

export interface ApplicationStatus {
	id: number;
	slug: string;
	name: string;
	kind: ApplicationStatusKind;
	color: string;
	displayOrder: number;
	isDefault: boolean;
	createdAt: string;
}

export interface ApplicationStatusUpsertPayload {
	slug?: string;
	name: string;
	kind: ApplicationStatusKind;
	color?: string;
	displayOrder?: number;
	isDefault?: boolean;
}

export interface ApplicationStatusEvent {
	id: number;
	applicationId: number;
	fromStatusId?: number;
	toStatusId: number;
	actorId?: number;
	note: string;
	createdAt: string;
	fromStatusName?: string;
	toStatusName?: string;
	actorEmail?: string;
}

export interface Application {
	id: number;
	jobId: number;
	jobTitle?: string;
	name: string;
	email: string;
	message: string;
	createdAt: string;
	statusId?: number;
	statusUpdatedAt?: string;
	statusUpdatedBy?: number;
	status?: ApplicationStatus;
	files?: ApplicationFile[];
	history?: ApplicationStatusEvent[];
}

export type ApplicationSort = 'newest' | 'oldest';

export interface ApplicationListFilter {
	jobId?: number;
	statusId?: number;
	sort?: ApplicationSort;
	limit?: number;
}
