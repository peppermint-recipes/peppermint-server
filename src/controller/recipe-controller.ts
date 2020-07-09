import { Request, Response, NextFunction } from "express";

import { Recipe } from "../entity/recipe";
import { RecipeService } from "../services/recipe-service";

export class RecipeController {
    public static async getAllRecipes(
        request: Request,
        response: Response
    ): Promise<void> {
        try {
            const recipes: Recipe[] = await RecipeService.getAllRecipes();
            response.status(200).send({ recipes });
        } catch (error) {
            response.status(500).send({ status: "Could not load recipes" });
        }
    }

    public static async createRecipe(
        request: Request,
        response: Response,
        next: NextFunction
    ): Promise<void> {
        const name = request.body.name;
        const servings = request.body.servings;
        const activeTime = request.body.activeTime;
        const totalTime = request.body.totalTime;
        const ingredients = request.body.ingredients;
        const instructions = request.body.instructions;

        if (
            !name ||
            servings === undefined ||
            activeTime === undefined ||
            totalTime === undefined ||
            ingredients === undefined ||
            instructions === undefined
        ) {
            response.status(400).send({ status: "Arguments missing" });

            return;
        }

        try {
            const createdRecipe = await RecipeService.createRecipe(
                request.body
            );

            response.send({ status: "ok", recipe: createdRecipe });
        } catch (error) {
            response.status(500).send({ status: error.message });
        }
    }

    public static async getRecipeByID(
        request: Request,
        response: Response,
        next: NextFunction
    ): Promise<void> {
        const recipeID = request.params.recipeID;

        if (!recipeID) {
            response.status(404).send({ status: "Not found" });

            return;
        }

        try {
            const recipe = await RecipeService.getRecipeByID({ recipeID });
            response.send({ status: "ok", recipe });
        } catch (error) {
            response.status(404).send({ status: "No recipe found" });
        }
    }

    public static async deleteRecipeByID(
        request: Request,
        response: Response,
        next: NextFunction
    ): Promise<void> {
        const recipeID = request.params.recipeID;

        if (!recipeID) {
            response.status(404).send({ status: "Not recipe found" });

            return;
        }

        try {
            const recipe = await RecipeService.deleteRecipeByID({
                recipeID,
            });
            response.send({ status: "ok", recipe });
        } catch (error) {
            response.status(404).send({ status: error.message });
        }
    }

    public static async updateRecipe(
        request: Request,
        response: Response,
        next: NextFunction
    ): Promise<void> {
        const recipeID = request.params.recipeID;

        if (recipeID === undefined) {
            response.status(404).send({ status: "No recipe found" });

            return;
        }
        try {
            const createdRecipe = await RecipeService.updateRecipeByID({
                recipeID,
                ...request.body,
            });

            response.send({ status: "ok", recipe: createdRecipe });
        } catch (error) {
            response.status(409).send({ status: error.message });
        }
    }
}
