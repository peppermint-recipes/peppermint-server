import {
    Controller,
    Get,
    Post,
    Delete,
    Body,
    HttpException,
    HttpStatus,
    Param,
    Put,
} from "@nestjs/common";
import { IsNotEmpty, IsString, IsNumberString } from "class-validator";
import { InjectRepository } from "@nestjs/typeorm";
import { Repository } from "typeorm";

import { RecipeService } from "./recipe-service";
import { Recipe } from "./recipe-model";

class RecipeDto {
    @IsNotEmpty()
    @IsString()
    public readonly name: string;

    @IsNotEmpty()
    @IsNumberString()
    public readonly servings: string;

    @IsNotEmpty()
    @IsNumberString()
    public readonly activeTime: string;

    @IsNotEmpty()
    @IsNumberString()
    public readonly totalTime: string;

    @IsNotEmpty()
    @IsString()
    public readonly ingredients: string;

    @IsNotEmpty()
    @IsString()
    public readonly instructions: string;
}

@Controller()
export class RecipeController {
    constructor(
        private readonly recipeService: RecipeService,
        @InjectRepository(Recipe)
        private readonly recipeRepository: Repository<Recipe>
    ) {}

    @Get("api/recipes/")
    async getRecipes() {
        try {
            return {
                recipes: await this.recipeRepository.find(),
            };
        } catch (err) {
            throw new HttpException(
                { status: HttpStatus.INTERNAL_SERVER_ERROR },
                500
            );
        }
    }

    @Post("api/recipes/")
    async createProduct(@Body() recipeDto: RecipeDto) {
        const recipe = new Recipe({
            name: recipeDto.name,
            servings: recipeDto.servings,
            activeTime: recipeDto.activeTime,
            totalTime: recipeDto.totalTime,
            ingredients: recipeDto.ingredients,
            instructions: recipeDto.instructions,
        });

        try {
            const createdRecipe = await this.recipeRepository.save(recipe);

            return createdRecipe;
        } catch (error) {
            throw new HttpException(
                { status: HttpStatus.INTERNAL_SERVER_ERROR },
                500
            );
        }
    }

    @Get("api/recipes/:id")
    async findRecipe(@Param("id") id: string) {
        try {
            const foundRecipe = await this.recipeRepository.findOneOrFail(id);

            return foundRecipe;
        } catch (error) {
            throw new HttpException(
                { status: HttpStatus.INTERNAL_SERVER_ERROR },
                500
            );
        }
    }

    @Put("api/recipes/:id")
    async updateRecipe(@Param("id") id: string, @Body() recipeDto: RecipeDto) {
        const recipe = new Recipe({
            name: recipeDto.name,
            servings: recipeDto.servings,
            activeTime: recipeDto.activeTime,
            totalTime: recipeDto.totalTime,
            ingredients: recipeDto.ingredients,
            instructions: recipeDto.instructions,
        });

        try {
            recipe.updatedAt = String(Date.now());
            const createdRecipe = await this.recipeRepository.save(recipe);

            return createdRecipe;
        } catch (error) {
            throw new HttpException(
                { status: HttpStatus.INTERNAL_SERVER_ERROR },
                500
            );
        }
    }

    @Delete("api/recipes/:id")
    async removeRecipe(@Param("id") id: string) {
        try {
            const deletedRecipe = await this.recipeRepository.delete(id);

            return deletedRecipe;
        } catch (error) {
            throw new HttpException(
                { status: HttpStatus.INTERNAL_SERVER_ERROR },
                500
            );
        }
    }
}
