import { Module } from "@nestjs/common";
import { TypeOrmModule } from "@nestjs/typeorm";

import { RecipeController } from "./recipe-controller";
import { RecipeService } from "./recipe-service";
import { Recipe } from "./recipe-model";

@Module({
    imports: [TypeOrmModule.forFeature([Recipe])],
    controllers: [RecipeController],
    providers: [RecipeService],
})
export class RecipeModule {}
