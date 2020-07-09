/** Package imports */
import { Router } from "express";
import { recipeRouter } from "./recipe-router";

/** Variables */
export const globalRouter: Router = Router({ mergeParams: true });

/** Routes */
globalRouter.use("/recipes", recipeRouter);
