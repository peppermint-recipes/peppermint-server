import { getRepository } from "typeorm";
import { Recipe } from "../entity/recipe";

export class RecipeService {
    private static getRecipeRepository() {
        return getRepository(Recipe);
    }

    public static async getAllRecipes() {
        const recipeRepository = this.getRecipeRepository();

        return await recipeRepository.find();
    }

    public static async createRecipe(options: {
        name: string;
        servings: string;
        activeTime: string;
        totalTime: string;
        ingredients: string;
        instructions: string;
    }): Promise<Recipe> {
        const recipeRepository = this.getRecipeRepository();

        try {
            const recipe = new Recipe({
                name: options.name,
                servings: options.servings,
                activeTime: options.activeTime,
                totalTime: options.totalTime,
                ingredients: options.ingredients,
                instructions: options.instructions,
            });

            return await recipeRepository.save(recipe);
        } catch (error) {
            throw new Error("Could not create recipe.");
        }
    }

    public static async updateRecipe(recipe: Recipe): Promise<Recipe> {
        const recipeRepository = this.getRecipeRepository();
        recipe.updatedAt = String(Date.now());

        return await recipeRepository.save(recipe);
    }

    public static async getRecipeByID(options: {
        recipeID: string;
    }): Promise<Recipe> {
        const recipeRepository = this.getRecipeRepository();
        try {
            const recipe = await recipeRepository.findOneOrFail(
                options.recipeID
            );

            return recipe;
        } catch (error) {
            throw new Error("Could not find recipe.");
        }
    }

    public static async deleteRecipeByID(options: { recipeID: string }) {
        const recipeRepository = this.getRecipeRepository();
        try {
            const recipe = await recipeRepository.findOneOrFail(
                options.recipeID
            );

            await recipeRepository.delete(recipe);
        } catch (error) {
            throw new Error("Could not delete recipe.");
        }
    }

    public static async updateRecipeByID(options: {
        recipeID?: string;
        name?: string;
        servings?: string;
        activeTime?: string;
        totalTime?: string;
        ingredients?: string;
        instructions?: string;
    }): Promise<Recipe> {
        const recipeRepository = this.getRecipeRepository();
        let recipe: Recipe;

        try {
            recipe = await recipeRepository.findOneOrFail(options.recipeID);
        } catch (error) {
            throw new Error("Could not find recipe.");
        }

        try {
            recipe.name = options.name ? options.name : recipe.name;
            recipe.servings = options.servings
                ? options.servings
                : recipe.servings;
            recipe.activeTime = options.activeTime
                ? options.activeTime
                : recipe.activeTime;
            recipe.totalTime = options.totalTime
                ? options.totalTime
                : recipe.totalTime;
            recipe.ingredients = options.ingredients
                ? options.ingredients
                : recipe.ingredients;
            recipe.instructions = options.instructions
                ? options.instructions
                : recipe.instructions;

            recipe.updatedAt = String(Date.now());

            return await recipeRepository.save(recipe);
        } catch (error) {
            throw new Error("Could not update recipe.");
        }
    }
}
