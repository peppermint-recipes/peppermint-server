import { Entity, PrimaryGeneratedColumn, Column } from "typeorm";

@Entity()
export class Recipe {
    @PrimaryGeneratedColumn()
    id: number;

    @Column()
    name: string;

    @Column()
    servings: string;

    @Column()
    activeTime: string;

    @Column()
    totalTime: string;

    @Column()
    ingredients: string;

    @Column()
    instructions: string;

    @Column()
    createdAt: string;

    @Column()
    updatedAt: string;

    public constructor(options?: {
        name: string;
        servings: string;
        activeTime: string;
        totalTime: string;
        ingredients: string;
        instructions: string;
    }) {
        if (!options) {
            return;
        }
        this.name = options.name;
        this.servings = options.servings;
        this.activeTime = options.activeTime;
        this.totalTime = options.totalTime;
        this.ingredients = options.ingredients;
        this.instructions = options.instructions;
        this.createdAt = String(Date.now());
        this.updatedAt = String(Date.now());
    }
}
