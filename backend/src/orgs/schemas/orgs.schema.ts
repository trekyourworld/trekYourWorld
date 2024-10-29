import { Prop, Schema, SchemaFactory } from '@nestjs/mongoose';
import { Document, HydratedDocument } from 'mongoose';

@Schema()
export class Organisations extends Document {
    @Prop()
    name: string;
    @Prop()
    label: string;
}

export type OrganisationDocument = HydratedDocument<Organisations>;

export const OrganisationSchema = SchemaFactory.createForClass(Organisations);
