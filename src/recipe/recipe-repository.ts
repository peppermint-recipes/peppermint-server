import { Injectable } from "@nestjs/common";
import { EntityRepository, Repository } from "typeorm";

import { Recipe } from "./recipe-model";

@Injectable()
@EntityRepository(Recipe)
export class RecipeRepository extends Repository<Recipe> {}
