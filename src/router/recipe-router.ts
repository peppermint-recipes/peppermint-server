import { Router } from "express";
import { RecipeController } from "../controller/recipe-controller";

export const recipeRouter: Router = Router({ mergeParams: true });

recipeRouter.get("/", RecipeController.getAllRecipes);
recipeRouter.post("/", RecipeController.createRecipe);
recipeRouter.get("/:recipeID", RecipeController.getRecipeByID);
recipeRouter.delete("/:recipeID", RecipeController.deleteRecipeByID);
recipeRouter.patch("/:recipeID", RecipeController.updateRecipe);
